<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class CreateAuditsTable extends Migration
{
    public function up()
    {
        $morphPrefix = config('audit.user.morph_prefix', 'user');

        Schema::create('audits', function (Blueprint $table) use ($morphPrefix) {
            $table->uuid('id')->primary();
            $table->string($morphPrefix.'_type')->nullable();
            $table->unsignedBigInteger($morphPrefix.'_id')->nullable();
            $table->string('event');
            $table->nullableUuidMorphs('auditable');
            $table->text('old_values')->nullable();
            $table->text('new_values')->nullable();
            $table->text('url')->nullable();
            $table->ipAddress('ip_address')->nullable();
            $table->string('user_agent', 1023)->nullable();
            $table->string('tags')->nullable();
            $table->timestamps();

            $table->index([$morphPrefix.'_id', $morphPrefix.'_type']);
        });
    }

    public function down()
    {
        Schema::dropIfExists('audits');
    }
}
