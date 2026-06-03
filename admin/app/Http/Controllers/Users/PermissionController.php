<?php

namespace App\Http\Controllers\Users;

use App\Http\Controllers\Controller;
use App\Http\Requests\PermissionRequest;
use App\Models\Permission\Permission;
use App\Utilities\UserUtilities;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class PermissionController extends Controller
{
    public function __construct(
        private readonly UserUtilities $userUtilities
    ) {}

    public function index(Request $request): Response
    {
        $userSettings = $this->userUtilities->syncUserSettings($request);

        $perPage = $request->query('perPage', $userSettings['settings_permissions_per_page'] ?? 10);
        $page = $request->query('page', 1);
        $search = request()->query('search', '');
        $filters = json_decode(request()->query('filters', '[]'), true);

        $perPage = is_numeric($perPage) ? (int) $perPage : 10;
        $page = is_numeric($page) ? (int) $page : 1;

        $uniqueNames = Permission::distinct()->pluck('name')->toArray();
        $uniqueRoutes = Permission::distinct()->pluck('route')->toArray();

        $query = Permission::query();

        if ($search) {
            $query->where(function ($query) use ($search) {
                $query->orWhere('name', 'ilike', '%' . $search . '%')
                    ->orWhere('route', 'ilike', '%' . $search . '%')
                    ->orWhere('description', 'ilike', '%' . $search . '%')
                ;
            });
        }

        if ($filters) {
            foreach ($filters as $key =>$values) {
                if (empty($values)) {
                    continue;
                }
                $query->whereIn($key, $values);
            }
        }

        $permissions = $query->paginate($perPage, ['*'], 'page', $page);

        return Inertia::render('permissions/index', [
            'data' => $permissions,
            'dataFilters' => ['name' => $uniqueNames, 'route' => $uniqueRoutes],
            'search' => $search,
            'activeFilters' => $filters,
        ]);
    }

    public function create(): Response
    {
        return Inertia::render('permissions/create');
    }

    public function store(PermissionRequest $request)
    {
        Permission::create([
            'name' => $request->name,
            'route' => $request->route,
            'description' => $request->description,
        ]);

        return redirect()->route('permissions.index')->with('success', 'Доступ успішно створено');
    }

    public function edit(Permission $permission): Response
    {
        return Inertia::render('permissions/edit', [
            'permission' => $permission,
        ]);
    }

    public function update(PermissionRequest $request, Permission $permission)
    {
        $permission->update([
            'name' => $request->name,
            'route' => $request->route,
            'description' => $request->description,
        ]);

        return redirect()->route('permissions.index')->with('success', 'Доступ успішно оновлено');
    }

    public function destroy(Permission $permission)
    {
        $permission->delete();

        return redirect()
            ->route('permissions.index')
            ->with('success', 'Доступ успішно видалено')
        ;
    }
}
