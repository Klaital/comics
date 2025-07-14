<?php

namespace App\Command;

use App\Entity\RssItem;
use App\Repository\RssItemRepository;
use App\Repository\SubscriptionRepository;
//use Composer\Console\Input\InputArgument;
use Doctrine\ORM\EntityManagerInterface;
use Psr\Log\LoggerInterface;
use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Console\Style\SymfonyStyle;
use Symfony\Component\Console\Input\InputArgument;

#[AsCommand(
    name: 'rss:unread',
    description: 'Mark RSS Items as unread for a single subscription',
    hidden: false,
)]
class UnreadRssFeeds extends Command
{
    private SubscriptionRepository $subscriptionRepository;
    private RSSItemRepository $rssItemRepository;
    private EntityManagerInterface $entityManager;
    private LoggerInterface $logger;

    public function __construct(
        EntityManagerInterface $entityManager,
        SubscriptionRepository $subscriptionRepository,
        RssItemRepository $rssItemRepository,
        LoggerInterface $logger,
    ) {
        parent::__construct();
        $this->entityManager = $entityManager;
        $this->subscriptionRepository = $subscriptionRepository;
        $this->rssItemRepository = $rssItemRepository;
        $this->logger = $logger;
    }

    protected function configure(): void
    {
        $this->addArgument("sub-id", InputArgument::REQUIRED, "Subscription ID to mark as unread");
    }

//    public function __invoke(#[Argument('ID of the subscription to reset')] int $usbId, InputInterface $input, OutputInterface $output)
//    {
//        $subscription = $this->subscriptionRepository->find($usbId);
//        if (!$subscription) {
//            $this->logger->error('Subscription not found');
//            return Command::INVALID;
//        }
//
//    }
    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $io = new SymfonyStyle($input, $output);
        $io->title('Marking RSS feed as Unread');
        $subId = $input->getArgument('sub-id');
        $io->writeln("Got Sub ID: " . $subId);
        if (!$subId) {
            $io->error('Sub ID not specified');
            return Command::INVALID;
        }
        $s = $this->subscriptionRepository->find($subId);
        if (!$s) {
            $io->error("Subscription not found");
            return Command::INVALID;
        }
        $output->writeln("Loaded subscription data for " . $s->getTitle());

        $s->markAllRssAsUnread($this->entityManager);
        $io->success("Subscription marked as unread");
        return Command::SUCCESS;
    }
}
