<?php

namespace App\MessageHandler;

use App\Message\ItemReadMessage;
use App\Repository\RssItemRepository;
use Doctrine\ORM\EntityManagerInterface;
use Symfony\Component\Messenger\Attribute\AsMessageHandler;

#[AsMessageHandler]
class ItemReadMessageHandler
{
    public function __construct(
        private EntityManagerInterface $entityManager,
        private RssItemRepository $rssItemRepository,
    ) {

    }

    public function __invoke(ItemReadMessage $message): void
    {
        $item = $this->rssItemRepository->find($message->getItemId());
        if  (!$item) {
            return;
        }

        $item->setReadAt($message->getReadAt());
        $this->entityManager->flush();
    }
}
