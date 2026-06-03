<?php

namespace App\Events;

use Illuminate\Broadcasting\Channel;
use Illuminate\Broadcasting\InteractsWithSockets;
use Illuminate\Broadcasting\PresenceChannel;
use Illuminate\Broadcasting\PrivateChannel;
use Illuminate\Contracts\Broadcasting\ShouldBroadcast;
use Illuminate\Foundation\Events\Dispatchable;
use Illuminate\Queue\SerializesModels;

class UserActivityEvent implements ShouldBroadcast
{
    use Dispatchable, InteractsWithSockets, SerializesModels;

    public $userId;
    public string $userName;
    public $action;
    public string $url;
    public string $timestamp;

    /**
     * Create a new event instance.
     */
    public function __construct($user, $action, $url)
    {
        $this->userId = $user->id;
        $this->userName = $user->name;
        $this->action = $action;
        $this->url = $url;
        $this->timestamp = now()->toIso8601String();
    }

    /**
     * Get the channels the event should broadcast on.
     *
     * @return array<int, Channel>
     */
    public function broadcastOn(): array
    {
        return [
            new Channel('admin'),
        ];
    }

    public function broadcastWith(): array
    {
        return [
            'userId' => $this->userId,
            'userName' => $this->userName,
            'action' => $this->action,
            'url' => $this->url,
            'timestamp' => $this->timestamp,
        ];
    }

    public function broadcastAs(): string
    {
        return 'user.activity';
    }
}
