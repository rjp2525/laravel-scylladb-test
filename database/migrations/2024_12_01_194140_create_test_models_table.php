<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('test_models', function (Blueprint $table) {
            $table->uuid('id')->primary();
            $table->string('field_one')->nullable();
            $table->string('field_two')->nullable();
            $table->string('field_three')->nullable();
            $table->timestamps();
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('test_models');
    }
};
