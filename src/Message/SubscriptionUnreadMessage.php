<?php

namespace App\Message;

class SubscriptionUnreadMessage
{
    public function __construct(private int $sub_id)
    {

    }

    public function getSubId(): int
    {
        return $this->sub_id;
    }
}
