<?php

namespace App\Controller;

use App\Entity\RssItem;
use App\Entity\Subscription;
use App\Form\SubscriptionType;
use App\Message\ItemReadMessage;
use Doctrine\ORM\EntityManagerInterface;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\Messenger\MessageBusInterface;
use Symfony\Component\Routing\Attribute\Route;

final class SubscriptionController extends AbstractController
{
    public function __construct(
        private EntityManagerInterface $entityManager,
        private MessageBusInterface $bus,
    ){

    }
    #[Route('/subs', name: 'app_subscription')]
    public function index(): Response
    {
        $subs = $this->entityManager->getRepository(Subscription::class)->findAll();
        return $this->render('subscription/index.html.twig', [
            'controller_name' => 'SubscriptionController',
            'subs' => $subs,
        ]);
    }

    #[Route('/subs/{subId}', name: 'app_subscription_update')]
    public function update(Request $request): Response
    {
        $sub = $this->entityManager->getRepository(Subscription::class)->find($request->get('subId'));
        return $this->render('subscription/update.html.twig', [
            'sub' => $sub,
        ]);
    }

    #[Route('/api/subs', name: 'api_subscription')]
    public function listSubs(): Response
    {
        $subs = $this->entityManager->getRepository(Subscription::class)->findAll();
        return $this->json($subs);
    }

    #[Route('/subs/new', methods: ['GET', 'POST'])]
    public function createSubscription(Request $request): Response
    {
        $s = new Subscription();
        $form = $this->createForm(SubscriptionType::class, $s);

        $form->handleRequest($request);

        if ($form->isSubmitted() && $form->isValid()) {
            $this->entityManager->persist($s);
            $entityManager->flush();
            $this->addFlash('success', 'Subscription created');
            return $this->redirectToRoute('app_subscription');
        }
        return $this->render('subscription/new.html.twig', [
            'form' => $form->createView(),
        ]);
    }

    #[Route('/api/subs/{subId}/rss/readall', name: 'api_subscription_readall', methods: ['POST'])]
    public function readAll(int $subId): Response
    {
        $s = $this->entityManager->getRepository(Subscription::class)->find($subId);
        if (!$s)
        {
            return $this->json(['error' => 'Subscription not found'], Response::HTTP_NOT_FOUND);
        }
        $rssItems = $s->getRssItems();
        $now = new \DateTime();
        foreach($rssItems as $rssItem) {
            $this->bus->dispatch(new ItemReadMessage($rssItem->getId(), $now));
            $rssItem->setReadAt(new \DateTime());
        }
        $this->entityManager->flush();
        return new JsonResponse(null, 202);
    }

    #[Route('/api/subs/{subId}/rss/{itemId}/readbefore', name: 'api_subscription_readbefore', methods: ['POST'])]
    public function readAllBefore(int $subId, int $itemId): Response
    {
        $s = $this->entityManager->getRepository(Subscription::class)->find($subId);
        if (!$s)
        {
            return $this->json(['error' => 'Subscription not found'], Response::HTTP_NOT_FOUND);
        }
        $item = $this->entityManager->getRepository(RssItem::class)->find($itemId);
        $conn =  $this->entityManager->getConnection();
        $sql = 'UPDATE rss_item SET read_at = CURRENT_TIMESTAMP WHERE subscription_id = :sub_id AND read_at IS NULL AND pub_date <= :pub_date';
        $res = $conn->executeQuery($sql, [
            'sub_id' => $s->getId(),
            'pub_date' => $item->getPubDate()->format('Y-m-d H:i:s'),
        ]);
        return new JsonResponse(['update'=>$res->rowCount()], 200);
    }

    #[Route('/api/subs/{subId}/unread', name: 'api_subscription_unread', methods: ['GET'])]
    public function unread(int $subId): Response
    {
        $s = $this->entityManager->getRepository(Subscription::class)->find($subId);
        if (!$s)
        {
            return $this->json(['error' => 'Subscription not found'], Response::HTTP_NOT_FOUND);
        }
        $rssItems = $s->getUnreadRssItems();
        return $this->json($rssItems, 200);
    }
}
