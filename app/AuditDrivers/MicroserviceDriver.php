<?php

namespace App\AuditDrivers;

use Illuminate\Support\Facades\Http;
use Illuminate\Support\Facades\Log;
use OwenIt\Auditing\Contracts\Audit;
use OwenIt\Auditing\Contracts\Auditable;
use OwenIt\Auditing\Contracts\AuditDriver;

class MicroserviceDriver implements AuditDriver
{
    protected $baseUrl;

    public function __construct()
    {
        $this->baseUrl = config('audit.drivers.microservice.base_url', 'http://localhost:8080');
    }

    public function audit(Auditable $model): Audit
    {
        $auditData = $this->transformAuditData($model->toAudit());

        $response = Http::post("{$this->baseUrl}/audit", $auditData);

        if (! $response->successful()) {
            Log::error('Failed to send audit log to Go microservice', [
                'status' => $response->status(),
                'body' => $response->body(),
            ]);

            throw new \Exception('Failed to send audit log to Go microservice');
        }

        return $this->createAuditInstance($auditData);
    }

    public function prune(Auditable $model): bool
    {
        $threshold = $model->getAuditThreshold();

        if ($threshold > 0) {
            $response = Http::delete("{$this->baseUrl}/audit/prune", [
                'auditable_id' => $model->getKey(),
                'auditable_type' => $model->getMorphClass(),
                'threshold' => $threshold,
            ]);

            if (! $response->successful()) {
                Log::error('Failed to prune audit logs', [
                    'status' => $response->status(),
                    'body' => $response->body(),
                ]);

                return false;
            }

            return true;
        }

        return false;
    }

    public function fetchAuditsByTenantId(string $tenantId): array
    {
        $response = Http::get("{$this->baseUrl}/audit", ['tenant_id' => $tenantId]);

        if (! $response->successful()) {
            Log::error('Failed to fetch audit logs by tenant ID', [
                'tenant_id' => $tenantId,
                'status' => $response->status(),
                'body' => $response->body(),
            ]);

            throw new \Exception('Failed to fetch audit logs from Go microservice');
        }

        return $response->json();
    }

    public function fetchAll(): array
    {
        $response = Http::get("{$this->baseUrl}/audits");

        if (! $response->successful()) {
            Log::error('Failed to fetch all audit logs', [
                'status' => $response->status(),
                'body' => $response->body(),
            ]);

            throw new \Exception('Failed to fetch all audit logs');
        }

        return $response->json();
    }

    protected function transformAuditData(array $data): array
    {
        return collect(array_merge($data, [
            'created_at' => $this->formatTimestamp(data_get($data, 'created_at', now())),
            'updated_at' => $this->formatTimestamp(data_get($data, 'updated_at', now())),
            'old_values' => json_encode($data['old_values']),
            'new_values' => json_encode($data['new_values']),
        ]))->filter()->toArray();
    }

    protected function formatTimestamp($timestamp): string
    {
        return \Carbon\Carbon::parse($timestamp)->toIso8601String();
    }

    protected function createAuditInstance(array $data): Audit
    {
        $auditModel = config('audit.implementation', \OwenIt\Auditing\Models\Audit::class);

        return new $auditModel($data);
    }
}
