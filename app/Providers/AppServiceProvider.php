<?php

namespace App\Providers;

use App\Auditing\Drivers\ScyllaAuditDriver;
use Illuminate\Support\ServiceProvider;
use OwenIt\Auditing\Auditor;

class AppServiceProvider extends ServiceProvider
{
    public function register(): void {}

    public function boot(): void
    {
        $this->app->make(Auditor::class)->extend('microservice', function ($app) {
            return new ScyllaAuditDriver;
        });
    }
}
