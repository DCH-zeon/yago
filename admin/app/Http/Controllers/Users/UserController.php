<?php

namespace App\Http\Controllers\Users;

use App\Http\Controllers\Controller;
use App\Models\User;
use App\Utilities\UserUtilities;
use Spatie\Permission\Models\Role;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Hash;
use Inertia\Inertia;
use Inertia\Response;

class UserController extends Controller
{
    public function __construct(
        private readonly UserUtilities $userUtilities
    ) {}

    public function index(Request $request): Response
    {
        $userSettings = $this->userUtilities->syncUserSettings($request);
        $perPage = $request->query('perPage', $userSettings['settings_users_per_page'] ?? 10);
        $page = $request->query('page', 1);
        $search = request()->query('search', '');
        $filters = json_decode(request()->query('filters', '[]'), true);

        $perPage = is_numeric($perPage) ? (int) $perPage : 10;
        $page = is_numeric($page) ? (int) $page : 1;

        $query = User::query()->with('roles');

        if ($search) {
            $query->where(function ($query) use ($search) {
                $query->orWhere('name', 'ilike', '%' . $search . '%')
                    ->orWhere('email', 'ilike', '%' . $search . '%')
                ;
            });
        }

        $users = $query->paginate($perPage, ['*'], 'page', $page);

        return Inertia::render('users/index', [
            'data' => $users,
            'dataFilters' => [],
            'search' => $search,
            'activeFilters' => [],
        ]);
    }

    public function create(): Response
    {
        $roles = Role::all();

        return Inertia::render('users/create', [
            'roles' => $roles,
        ]);
    }

    public function store(Request $request)
    {
        $validated = $request->validate([
            'name' => 'required|string|max:255',
            'email' => 'required|email|unique:users',
            'password' => 'required|string|min:8|confirmed',
            'password_confirmation' => 'required|string|min:8',
            'roles' => 'array',
        ]);

        $user = User::create([
            'name' => $validated['name'],
            'email' => $validated['email'],
            'password' => Hash::make($validated['password']),
        ]);

        if (!empty($validated['roles'])) {
            $user->syncRoles($validated['roles']);
        }

        return redirect()->route('users.index')->with('success', 'Користувач успішно створений');
    }

    public function edit(User $user): Response
    {
        $roles = Role::all();
        $userRoles = $user->roles->pluck('id')->toArray();

        return Inertia::render('users/edit', [
            'user' => $user,
            'roles' => $roles,
            'userRoles' => $userRoles,
        ]);
    }

    public function update(Request $request, User $user)
    {
        $validated = $request->validate([
            'name' => 'required|string|max:255',
            'email' => 'required|email|unique:users,email,' . $user->id,
            'password' => 'nullable|string|min:8|confirmed',
            'roles' => 'array',
        ]);

        $user->update([
            'name' => $validated['name'],
            'email' => $validated['email'],
        ]);

        if (!empty($validated['password'])) {
            $user->update(['password' => Hash::make($validated['password'])]);
        }

        $user->syncRoles($validated['roles'] ?? []);

        return redirect()->route('users.index')->with('success', 'Користувач успішно оновлений');
    }

    public function destroy(User $user)
    {
        if ($user->id === auth()->id()) {
            return redirect()->route('users.index')->with('error', 'Ви не можете видалити себе');
        }

        $user->delete();

        return redirect()->route('users.index')->with('success', 'Користувач успішно видалений');
    }
}
