<?php

namespace App\Message;

class CheckRssMessage
{
    public function __construct(private int $sub_id, private bool $readall=false)
    {

    }

    public function getSubId(): int
    {
        return $this->sub_id;
    }

    public function getReadall(): bool
    {
        return $this->readall;
    }
}
