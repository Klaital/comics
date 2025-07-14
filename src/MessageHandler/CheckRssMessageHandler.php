<?php

namespace App\MessageHandler;

use App\Message\CheckRssMessage;
use Doctrine\ORM\EntityManagerInterface;
use Symfony\Component\Messenger\Attribute\AsMessageHandler;

#[AsMessageHandler]
class CheckRssMessageHandler
{
    public function __construct(
        private EntityManagerInterface $entityManager,
    ) {

    }

    public function __invoke(CheckRssMessage $message)
    {

    }
}
