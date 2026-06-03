<?php

namespace App\Utilities;

use Illuminate\Http\Request;

class UserUtilities
{
    public function syncUserSettings(Request $request): array
    {
        $user = $request->user();
        $settingsFromRequest = json_decode(rawurldecode($request->cookie('settings', '[]')), true);

        if ($user->settings['is_remember_settings'] ?? false) {
            $user->settings = array_merge($user->settings, $settingsFromRequest);
        }

        $user->save();

        return $user->settings;
    }
}
