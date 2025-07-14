<?php

namespace App\Message;

class ItemReadMessage
{
    public function __construct(private int $item_id, private \DateTime $read_at)
    {

    }

    public function getItemId(): int
    {
        return $this->item_id;
    }

    public function getReadAt(): \DateTime
    {
        return $this->read_at;
    }




}
