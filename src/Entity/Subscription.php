<?php

namespace App\Entity;

use App\Repository\RssItemRepository;
use App\Repository\SubscriptionRepository;
use Doctrine\Common\Collections\ArrayCollection;
use Doctrine\Common\Collections\Collection;
use Doctrine\DBAL\Types\Types;
use Doctrine\ORM\EntityManagerInterface;
use Doctrine\ORM\Mapping as ORM;
use Exception;
use Psr\Log\LoggerInterface;
use Symfony\Contracts\HttpClient\HttpClientInterface;
use Symfony\UX\Turbo\Attribute\Broadcast;

#[ORM\Entity(repositoryClass: SubscriptionRepository::class)]
#[Broadcast]
class Subscription
{
    function __construct()
    {
        $this->created_at = new \DateTime();
        $this->updated_at = new \DateTime();
        $this->read_at = new \DateTime();
        $this->rssItems = new ArrayCollection();
    }

    #[ORM\Id]
    #[ORM\GeneratedValue]
    #[ORM\Column]
    private ?int $id = null;

    #[ORM\Column(length: 255)]
    private ?string $title = null;

    #[ORM\Column(nullable: true)]
    private ?int $ordinal = null;

    #[ORM\Column(length: 255)]
    private ?string $homepage = null;

    #[ORM\Column(length: 255, nullable: true)]
    private ?string $first_comic_url = null;

    #[ORM\Column(length: 255, nullable: true)]
    private ?string $lastest_comic_url = null;

    #[ORM\Column(length: 255, nullable: true)]
    private ?string $rss_url = null;

    #[ORM\Column(nullable: false)]
    private ?bool $updates_mon = false;

    #[ORM\Column(nullable: false)]
    private ?bool $updates_tue = false;

    #[ORM\Column(nullable: false)]
    private ?bool $updates_wed = false;

    #[ORM\Column(nullable: false)]
    private ?bool $updates_thu = false;

    #[ORM\Column(nullable: false)]
    private ?bool $updates_fri = false;

    #[ORM\Column(nullable: false)]
    private ?bool $updates_sat = false;

    #[ORM\Column(nullable: false)]
    private ?bool $updates_sun = false;

    #[ORM\Column(nullable: false)]
    private ?bool $active = false;

    #[ORM\Column(nullable: false)]
    private ?bool $nsfw = false;

    #[ORM\Column(type: Types::DATETIMETZ_MUTABLE, nullable: false, options: ["default" => "CURRENT_TIMESTAMP"])]
    private ?\DateTime $created_at = null;

    #[ORM\Column(type: Types::DATETIMETZ_MUTABLE, nullable: false, options: ["default" => "CURRENT_TIMESTAMP"])]
    private ?\DateTime $updated_at = null;

    #[ORM\Column(type: Types::DATETIMETZ_MUTABLE, nullable: false, options: ["default" => "CURRENT_TIMESTAMP"])]
    private ?\DateTime $read_at = null;

    private int $rss_unread_count = 0;

    /**
     * @var Collection<int, RssItem>
     */
    #[ORM\OneToMany(targetEntity: RssItem::class, mappedBy: 'subscription', orphanRemoval: true)]
    private Collection $rssItems;

    /**
     * @return Collection<int, RssItem>
     * @throws Exception
     */
    public function getUnreadRssItems(): Collection
    {

        $unread = $this->rssItems->filter(fn(RssItem $item) => $item->getReadAt() === null);
        $it = $unread->getIterator();
        $it->uasort(fn(RssItem $a, RssItem $b) => $a->getPubDate() <=> $b->getPubDate());
        return new ArrayCollection(iterator_to_array($it));
    }

    public function getUnreadCount(): int
    {
        return $this->getUnreadRssItems()->count();
    }

    /**
     * Fetches the oldest unread RssItem for a subscription
     * @return ?RssItem
     */
    public function getOldestUnreadItem(): ?RssItem
    {
        $unread = $this->getUnreadRssItems();
        if ($unread->count() === 0) {
            return null;
        }
        return $unread->first();
    }

    public function getId(): ?int
    {
        return $this->id;
    }

    public function getTitle(): ?string
    {
        return $this->title;
    }

    public function setTitle(string $title): static
    {
        $this->title = $title;

        return $this;
    }

    public function getOrdinal(): ?int
    {
        return $this->ordinal;
    }

    public function setOrdinal(?int $ordinal): static
    {
        $this->ordinal = $ordinal;

        return $this;
    }

    public function getHomepage(): ?string
    {
        return $this->homepage;
    }

    public function setHomepage(string $homepage): static
    {
        $this->homepage = $homepage;

        return $this;
    }

    public function getFirstComicUrl(): ?string
    {
        return $this->first_comic_url;
    }

    public function setFirstComicUrl(?string $first_comic_url): static
    {
        $this->first_comic_url = $first_comic_url;

        return $this;
    }

    public function getLastestComicUrl(): ?string
    {
        return $this->lastest_comic_url;
    }

    public function setLastestComicUrl(?string $lastest_comic_url): static
    {
        $this->lastest_comic_url = $lastest_comic_url;

        return $this;
    }

    public function getRssUrl(): ?string
    {
        return $this->rss_url;
    }

    public function setRssUrl(?string $rss_url): static
    {
        $this->rss_url = $rss_url;

        return $this;
    }

    public function isUpdatesMon(): ?bool
    {
        return $this->updates_mon;
    }

    public function setUpdatesMon(bool $updates_mon): static
    {
        $this->updates_mon = $updates_mon;

        return $this;
    }

    public function isUpdatesTue(): ?bool
    {
        return $this->updates_tue;
    }

    public function setUpdatesTue(bool $updates_tue): static
    {
        $this->updates_tue = $updates_tue;

        return $this;
    }

    public function isUpdatesWed(): ?bool
    {
        return $this->updates_wed;
    }

    public function setUpdatesWed(bool $updates_wed): static
    {
        $this->updates_wed = $updates_wed;

        return $this;
    }

    public function isUpdatesThu(): ?bool
    {
        return $this->updates_thu;
    }

    public function setUpdatesThu(bool $updates_thu): static
    {
        $this->updates_thu = $updates_thu;

        return $this;
    }

    public function isUpdatesFri(): ?bool
    {
        return $this->updates_fri;
    }

    public function setUpdatesFri(bool $updates_fri): static
    {
        $this->updates_fri = $updates_fri;

        return $this;
    }

    public function isUpdatesSat(): ?bool
    {
        return $this->updates_sat;
    }

    public function setUpdatesSat(bool $updates_sat): static
    {
        $this->updates_sat = $updates_sat;

        return $this;
    }

    public function isUpdatesSun(): ?bool
    {
        return $this->updates_sun;
    }

    public function setUpdatesSun(bool $updates_sun): static
    {
        $this->updates_sun = $updates_sun;

        return $this;
    }

    public function isActive(): ?bool
    {
        return $this->active;
    }

    public function setActive(bool $active): static
    {
        $this->active = $active;

        return $this;
    }

    public function isNsfw(): ?bool
    {
        return $this->nsfw;
    }

    public function setNsfw(bool $nsfw): static
    {
        $this->nsfw = $nsfw;

        return $this;
    }

    public function getCreatedAt(): ?\DateTime
    {
        return $this->created_at;
    }

    public function setCreatedAt(\DateTime $created_at): static
    {
        $this->created_at = $created_at;

        return $this;
    }

    public function getUpdatedAt(): ?\DateTime
    {
        return $this->updated_at;
    }

    public function setUpdatedAt(\DateTime $updated_at): static
    {
        $this->updated_at = $updated_at;

        return $this;
    }

    public function getReadAt(): ?\DateTime
    {
        return $this->read_at;
    }

    public function setReadAt(\DateTime $read_at): static
    {
        $this->read_at = $read_at;

        return $this;
    }

    /**
     * @return Collection<int, RssItem>
     */
    public function getRssItems(): Collection
    {
        return $this->rssItems;
    }

    public function addRssItem(RssItem $rssItem): static
    {
        if (!$this->rssItems->contains($rssItem)) {
            $this->rssItems->add($rssItem);
            $rssItem->setSubscription($this);
        }

        return $this;
    }

    public function removeRssItem(RssItem $rssItem): static
    {
        if ($this->rssItems->removeElement($rssItem)) {
            // set the owning side to null (unless already changed)
            if ($rssItem->getSubscription() === $this) {
                $rssItem->setSubscription(null);
            }
        }
        return $this;
    }

    public function markAllRssAsRead(
        EntityManagerInterface $entityManager,
    ): void
    {
        $conn = $entityManager->getConnection();
        $sql = 'UPDATE rss_item SET read_at = :read_at WHERE subscription_id  = :sub_id AND read_at IS NULL';
        $res = $conn->executeQuery($sql, ['read_at' => new \DateTime(), 'sub_id' => $this->getId()]);
    }

    public function markAllRssAsUnread(
        EntityManagerInterface $entityManager,
    ): void
    {
        $conn = $entityManager->getConnection();
        $sql = 'UPDATE rss_item SET read_at = NULL WHERE subscription_id  = :sub_id AND read_at IS NOT NULL';
        $res = $conn->executeQuery($sql, ['sub_id' => $this->getId()]);
    }


    public function checkRssFeed(
        EntityManagerInterface $entityManager,
        HttpClientInterface $httpClient,
        RssItemRepository $rssItemRepository,
        LoggerInterface $logger): void
    {
        // TODO: copy in the code from ReadRssFeeds to here
        // Fetch the RSS feed
        $resp = $httpClient->request('GET', $this->rss_url);
        if ($resp->getStatusCode() !== 200) {
            throw new \RuntimeException("Failed to fetch ($this->rss_url): HTTP {$resp->getStatusCode()}");
        }

        $rssContent = $resp->getContent(true);
        $rssData = @simplexml_load_string($rssContent);
        if (!$rssData) {
            throw new \RuntimeException("Failed to parse rss feed " .  $this->rss_url);
        }
        foreach ($rssData->channel->item as $item) {
            $newItem = new RssItem();
            $title = (string) $item->title;
            $description = (string) $item->description;
            $link = (string) $item->link;
            $guid = (string) ($item->guid ?? $item->link);
            $pubDateRaw = (string) $item->pubDate;
            $pubDate = new \DateTime();
            try {
                $pubDate = new \DateTime($pubDateRaw);
            } catch (Exception $e) {
                $logger->warning("Failed to parse pub date '{pubDateRaw}'", ['pubDateRaw' => $pubDateRaw]);
            }

            try {
                $item = $rssItemRepository->findOneBy(array('guid' => $guid));
                // Try to fetch an existing item with the same GUID from the DB
                // If found, skip this item - we're only adding new items
                if ($item != null) {
                    print('.');
                    continue;
                }
            } catch (Exception $e) {
                $logger->warning("Failed to fetch item ($title): {$e->getMessage()}");
            }

            $newItem->setTitle($title);
            $newItem->setDescription($description);
            $newItem->setLink($link);
            $newItem->setPubDate($pubDate);
            $newItem->setGuid($guid);

            $logger->debug('New Item: ' . $newItem->getTitle());
            $this->addRssItem($newItem);
            $entityManager->persist($newItem);
            print('+');
        }
        $entityManager->flush();
    }
}
