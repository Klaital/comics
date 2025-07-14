<?php

namespace App\MessageHandler;

use App\Entity\Subscription;
use App\Message\SubscriptionReadallMessage;
use Doctrine\ORM\EntityManagerInterface;
use Doctrine\ORM\Query\ResultSetMapping;
use Symfony\Component\Messenger\Attribute\AsMessageHandler;

#[AsMessageHandler]
class SubscriptionReadallMessageHandler
{
    public function __construct(
        private EntityManagerInterface $entityManager,
    ) {

    }

    public function __invoke(SubscriptionReadallMessage $message)
    {
        $s = $this->entityManager->getRepository(Subscription::class)->find($message->getSubId());
        $s->markAllRssAsRead($this->entityManager);
    }
}
