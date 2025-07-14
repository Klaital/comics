<?php

namespace App\Command;

use App\Entity\RssItem;
use App\Repository\RssItemRepository;
use App\Repository\SubscriptionRepository;
use Doctrine\ORM\EntityManagerInterface;
use Psr\Log\LoggerInterface;
use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Console\Style\SymfonyStyle;
use Symfony\Contracts\HttpClient\HttpClientInterface;

#[AsCommand(
    name: 'rss:feeds',
    description: 'Read RSS Feeds for all subscriptions',
    hidden: false,
)]
class ReadRssFeeds extends Command
{
    private SubscriptionRepository $subscriptionRepository;
    private RSSItemRepository $rssItemRepository;
    private HttpClientInterface $httpClient;
    private EntityManagerInterface $entityManager;
    private LoggerInterface $logger;

    public function __construct(
        EntityManagerInterface $entityManager,
        SubscriptionRepository $subscriptionRepository,
        RssItemRepository $rssItemRepository,
        HttpClientInterface $httpClient,
        LoggerInterface $logger,
    ) {
        parent::__construct();
        $this->entityManager = $entityManager;
        $this->subscriptionRepository = $subscriptionRepository;
        $this->rssItemRepository = $rssItemRepository;
        $this->httpClient = $httpClient;
        $this->logger = $logger;
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $output->writeln('<info>Reading RSS Feeds</info>');
        $io = new SymfonyStyle($input, $output);
        $subs = $this->subscriptionRepository->findAll();
        $count = count($subs);
        if ($count <= 0) {
            $io->warning('No subscriptions found');
            return Command::SUCCESS;
        }
        $io->title("Fetching $count RSS Feeds");

        foreach ($subs as $sub) {
            $io->section($sub->getTitle());
            try {
                $sub->checkRssFeed(
                    $this->entityManager, $this->httpClient,
                    $this->rssItemRepository, $this->logger);
            } catch (\RuntimeException $e) {
                $io->error($e->getMessage());
            }

        }
        $io->success("RSS updates complete");
        return Command::SUCCESS;
    }
}
