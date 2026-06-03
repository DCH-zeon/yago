<?php

namespace App\Http\Middleware;

use App\Events\UserActivityEvent;
use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

class TrackUserActivity
{
    /**
     * Handle an incoming request.
     *
     * @param  Closure(Request): (Response)  $next
     */
    public function handle(Request $request, Closure $next): Response
    {
        $response = $next($request);

        if ($request->user() && $request->isMethod('GET')) {
            $routeName = $request->route() ? ($request->route()->getName() ?? $request->path()) : $request->path();

            event(new UserActivityEvent(
                $request->user(),
                "Відвідав сторінку: $routeName",
                $request->fullUrl()
            ));
        }

        return $response;
    }
}
