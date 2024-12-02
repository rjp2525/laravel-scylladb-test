<?php

namespace App\Models;

use App\AuditDrivers\MicroserviceDriver;
use App\Casts\EncryptedArray;
use Illuminate\Database\Eloquent\Concerns\HasUuids;
use OwenIt\Auditing\Models\Audit as BaseAudit;

class Audit extends BaseAudit
{
    use HasUuids;

    protected $casts = [
        'old_values' => EncryptedArray::class,
        'new_values' => EncryptedArray::class,
    ];

    public static function fetchAllAudits(): array
    {
        $driver = app(MicroserviceDriver::class);

        if (method_exists($driver, 'fetchAll')) {
            return $driver->fetchAll();
        }

        throw new \Exception('The current audit driver does not support fetching all audits.');
    }
}
