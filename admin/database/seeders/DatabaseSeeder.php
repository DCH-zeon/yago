<?php

namespace Database\Seeders;

use App\Models\User;
use Illuminate\Database\Seeder;
use Spatie\Permission\Models\Role;
use Spatie\Permission\Models\Permission;

class DatabaseSeeder extends Seeder
{
    /**
     * Seed the application's database.
     */
    public function run(): void
    {
        // Створення прав доступу
        $permissions = [
            ['name' => 'view',      'route' => 'roles',         'description' => 'Перегляд ролей користувачів'],
            ['name' => 'create',    'route' => 'roles',         'description' => 'Створення ролей користувачів'],
            ['name' => 'edit',      'route' => 'roles',         'description' => 'Редагування ролей користувачів'],
            ['name' => 'delete',    'route' => 'roles',         'description' => 'Видалення ролей користувачів'],
            ['name' => 'view',      'route' => 'permissions',   'description' => 'Перегляд прав доступу'],
            ['name' => 'create',    'route' => 'permissions',   'description' => 'Створення прав доступу'],
            ['name' => 'edit',      'route' => 'permissions',   'description' => 'Редагування прав доступу'],
            ['name' => 'delete',    'route' => 'permissions',   'description' => 'Видалення прав доступу'],
            ['name' => 'view',      'route' => 'users',         'description' => 'Перегляд користувачів'],
            ['name' => 'create',    'route' => 'users',         'description' => 'Створення користувачів'],
            ['name' => 'edit',      'route' => 'users',         'description' => 'Редагування користувачів'],
            ['name' => 'delete',    'route' => 'users',         'description' => 'Видалення користувачів'],
        ];

        foreach ($permissions as $permission) {
            Permission::firstOrCreate(['name' => $permission['name'], 'route' => $permission['route']], ['description' => $permission['description']]);
        }

        // Створення ролі з усіма правами
        $adminRole = Role::firstOrCreate(['name' => 'admin'], ['description' => 'Адміністратор системи з усіма правами']);
        $adminRole->syncPermissions(Permission::all());

        // Створення користувача Admin та призначення йому ролі
        $admin = User::firstOrCreate(
            ['email' => config('auth.admin.login', 'admin@mail.com')],
            [
                'name' => 'Admin',
                'avatar' => 'https://github.com/shadcn.png',
                'password' => config('auth.admin.password', 'secret856'),
            ]
        );

        $admin->assignRole($adminRole);

        Role::firstOrCreate(['name' => 'Контент менеджер'], ['description' => 'Доступ до інструментів роботи над контентом']);
    }
}
