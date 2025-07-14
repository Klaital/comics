<?php

namespace App\Message;

class SubscriptionReadallMessage
{
    public function __construct(private int $sub_id, private \DateTime $read_at)
    {

    }

    public function getSubId(): int
    {
        return $this->sub_id;
    }

    public function getReadAt(): \DateTime
    {
        return $this->read_at;
    }
}
