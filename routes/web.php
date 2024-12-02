<?php

use App\Models\Audit;
use Illuminate\Support\Facades\Route;

Route::get('/', function () {
    return view('welcome');
});

Route::get('test', function () {
    return response()->json(Audit::fetchAllAudits());
});
