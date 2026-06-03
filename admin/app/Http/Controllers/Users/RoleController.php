<?php

namespace App\Http\Controllers\Users;

use App\Http\Controllers\Controller;
use App\Models\Permission\Permission;
use App\Models\Permission\Role;
use App\Utilities\UserUtilities;
use Illuminate\Support\Facades\Auth;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class RoleController extends Controller
{
    public function __construct(
        private readonly UserUtilities $userUtilities
    ) {}

    public function index(Request $request): Response
    {
        $userSettings = $this->userUtilities->syncUserSettings($request);
        $perPage = $request->query('perPage', $userSettings['settings_roles_per_page'] ?? 10);
        $page = $request->query('page', 1);
        $search = request()->query('search', '');
        $filters = json_decode(request()->query('filters', '[]'), true);

        $perPage = is_numeric($perPage) ? (int) $perPage : 10;
        $page = is_numeric($page) ? (int) $page : 1;

        $query = Role::query()->with('permissions');

        if ($search) {
            $query->where(function ($query) use ($search) {
                $query->orWhere('name', 'ilike', '%' . $search . '%')
                    ->orWhere('description', 'ilike', '%' . $search . '%')
                ;
            });
        }

        $roles = $query->paginate($perPage, ['*'], 'page', $page);

        return Inertia::render('roles/index', [
            'data' => $roles,
            'dataFilters' => [],
            'search' => $search,
            'activeFilters' => [],
        ]);
    }

    public function create(): Response
    {
        $permissions = Permission::all();

        return Inertia::render('roles/create', [
            'permissions' => $permissions,
        ]);
    }

    public function store(Request $request)
    {
        $validated = $request->validate(
            [
                'name' => 'required|string|unique:roles',
                'description' => 'nullable|string|max:255',
                'permissions' => 'array',
                'permissions.*' => 'exists:permissions,id',
            ],
            [
                'name.required' => 'Поле "Назва ролі" є обов\'язковим для заповнення.',
                'name.string' => 'Назва ролі має бути рядком.',
                'name.unique' => 'Така роль існує.',
                'description.string' => 'Опис має бути рядком.',
                'description.max' => 'Опис не може перевищувати 255 символів.',
                'permissions.array' => 'Список прав має бути масивом.',
                'permissions.*.exists' => 'Вибране право не існує.',
            ]
        );

        $role = Role::create(['name' => $validated['name'], 'description' => $validated['description']]);

        if (!empty($validated['permissions'])) {
            $role->permissions()->sync($validated['permissions']);
        }

        return redirect()->route('roles.index')->with('success', 'Роль успішно створена');
    }

    public function edit(Role $role): Response
    {
        $permissions = Permission::all();
        $rolePermissions = $role->permissions->pluck('id')->toArray();

        return Inertia::render('roles/edit', [
            'role' => $role,
            'permissions' => $permissions,
            'rolePermissions' => $rolePermissions,
        ]);
    }

    public function update(Request $request, Role $role)
    {
        $validated = $request->validate(
            [
                'name' => 'required|string|unique:roles,name,' . $role->id,
                'description' => 'nullable|string|max:255',
                'permissions' => 'array',
            ],
            [
                'name.required' => 'Поле "Назва ролі" є обов\'язковим для заповнення.',
                'name.string' => 'Назва ролі має бути рядком.',
                'name.unique' => 'Така роль існує.',
                'description.string' => 'Опис має бути рядком.',
                'description.max' => 'Опис не може перевищувати 255 символів.',
                'permissions.array' => 'Список прав має бути масивом.',
            ]
        );

        $role->update(['name' => $validated['name'], 'description' => $validated['description']]);

        $permissionIds = array_map('intval', $validated['permissions'] ?? []);
        $role->permissions()->sync($permissionIds);


        return redirect()->route('roles.index')->with('success', 'Роль успішно оновлено');
    }

    public function destroy(Role $role)
    {
        if ($role->name === 'admin') {
            return redirect()->route('roles.index')->with('error', 'Неможливо видалити роль адміністратора');
        }

        $role->delete();

        return redirect()->route('roles.index')->with('success', 'Роль успішно видалена');
    }
}
