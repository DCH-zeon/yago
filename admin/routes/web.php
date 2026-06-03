<?php

use App\Http\Controllers\Centrifugo\CentrifugoController;
use App\Http\Controllers\Users\PermissionController;
use App\Http\Controllers\Users\RoleController;
use App\Http\Controllers\Users\UserController;
use Illuminate\Support\Facades\Route;




//Route::middleware(['auth', 'auth.session'])->group(function () {
//    Route::inertia('/', 'welcome', ['canRegister' => false, 'canLogin' => false])->name('home');
//
//});


Route::middleware(['auth'])->group(function () {
    Route::get('/centrifugo/token', [CentrifugoController::class, 'getToken']);
    Route::get('/centrifugo/presence', [CentrifugoController::class, 'getPresence']);

    Route::inertia('/dashboard', 'dashboard')->name('dashboard');
    Route::inertia('/', 'dashboard')->name('home');
    Route::resource('roles', RoleController::class);
    Route::resource('permissions', PermissionController::class);
    Route::resource('users', UserController::class);
});

require __DIR__.'/settings.php';
