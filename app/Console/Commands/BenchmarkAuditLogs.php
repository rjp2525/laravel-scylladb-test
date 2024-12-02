<?php

namespace App\Console\Commands;

use App\AuditDrivers\MicroserviceDriver;
use App\Models\TestModel;
use Illuminate\Console\Command;
use OwenIt\Auditing\Models\Audit;

class BenchmarkAuditLogs extends Command
{
    protected $signature = 'benchmark:audit-logs {iterations=2500}';

    protected $description = 'Benchmark storing and retrieving audit logs in ScyllaDB vs MariaDB';

    public function handle()
    {
        $iterations = (int) $this->argument('iterations');

        // Benchmark storing logs
        $this->info('Starting audit log storage benchmark...');
        $startMariaDB = microtime(true);
        $this->benchmarkDB($iterations);
        $endMariaDB = microtime(true);

        $startScyllaDB = microtime(true);
        $this->benchmarkScyllaDB($iterations);
        $endScyllaDB = microtime(true);

        $this->info('MariaDB Store Time: '.($endMariaDB - $startMariaDB).' seconds');
        $this->info('ScyllaDB Store Time: '.($endScyllaDB - $startScyllaDB).' seconds');

        // Benchmark retrieval
        $this->info('Starting audit log retrieval benchmark...');
        $startMariaDB = microtime(true);
        $this->retrieveMariaDB();
        $endMariaDB = microtime(true);

        $startScyllaDB = microtime(true);
        $this->retrieveScyllaDB();
        $endScyllaDB = microtime(true);

        $this->info('MariaDB Retrieve Time: '.($endMariaDB - $startMariaDB).' seconds');
        $this->info('ScyllaDB Retrieve Time: '.($endScyllaDB - $startScyllaDB).' seconds');
    }

    protected function benchmarkDB($iterations, $driver = 'database')
    {
        config(['audit.driver' => $driver]);

        for ($i = 0; $i < $iterations; $i++) {
            $model = TestModel::create([
                'field_one' => 'Value '.$i,
                'field_two' => rand(1, 100),
                'field_three' => now()->toString(),
            ]);

            $m = TestModel::find($model->id);
            $m->update([
                'field_one' => 'Updated Value '.$i,
            ]);
        }
    }

    protected function benchmarkScyllaDB($iterations)
    {
        $this->benchmarkDB($iterations, MicroserviceDriver::class);
    }

    protected function retrieveMariaDB()
    {
        config(['audit.driver' => 'database']);

        Audit::all();
    }

    protected function retrieveScyllaDB()
    {
        config(['audit.driver' => MicroserviceDriver::class]);

        app(MicroserviceDriver::class)->fetchAll();
    }
}
