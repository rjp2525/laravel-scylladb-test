<?php

namespace App\Casts;

use Illuminate\Contracts\Database\Eloquent\CastsAttributes;
use Illuminate\Support\Facades\Crypt;

class EncryptedArray implements CastsAttributes
{
    public function get($model, string $key, $value, array $attributes)
    {
        if (is_null($value)) {
            return [];
        }

        return json_decode(Crypt::decryptString($value), true);
    }

    public function set($model, string $key, $value, array $attributes)
    {
        return Crypt::encryptString(json_encode($value));
    }
}
