<?php

namespace App\MessageHandler;

use App\Entity\Subscription;
use App\Message\SubscriptionUnreadMessage;
use Doctrine\ORM\EntityManagerInterface;
use Doctrine\ORM\Query\ResultSetMapping;
use Symfony\Component\Messenger\Attribute\AsMessageHandler;

#[AsMessageHandler]
class SubscriptionUnreadMessageHandler
{
    public function __construct(
        private EntityManagerInterface $entityManager,
    ) {

    }

    public function __invoke(SubscriptionUnreadMessage $message)
    {
        $s = $this->entityManager->getRepository(Subscription::class)->find($message->getSubId());
        $s->markAllRssAsUnread($this->entityManager);
    }
}
