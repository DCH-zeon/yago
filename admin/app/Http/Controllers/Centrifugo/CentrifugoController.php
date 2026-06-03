<?php

namespace App\Http\Controllers\Centrifugo;

use App\Http\Controllers\Controller;
use Illuminate\Http\Request;
use denis660\Centrifugo\Centrifugo;

class CentrifugoController extends Controller
{
    public function __construct(
        private readonly Centrifugo $centrifugo
    ) {}
    public function getToken(Request $request)
    {
        $user = $request->user();

        $token = $this->centrifugo->generateConnectionToken((string)$user->id, 3600 * 24, [
            'name' => $user->name,
            'avatar' => $user->avatar,
        ]);

        return response()->json(['token' => $token]);
    }

    public function getPresence()
    {
        $presence = $this->centrifugo->presence('admin');

        return response()->json($presence);
    }
}
