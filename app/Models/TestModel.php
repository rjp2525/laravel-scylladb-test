<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Model;
use OwenIt\Auditing\Auditable;
use OwenIt\Auditing\Contracts\Auditable as AuditableContract;

class TestModel extends Model implements AuditableContract
{
    use Auditable, HasUuids;

    protected $fillable = ['field_one', 'field_two', 'field_three'];
}
