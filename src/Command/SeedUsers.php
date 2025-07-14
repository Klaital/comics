<?php

namespace App\Command;

use App\Entity\User;
use App\Repository\UserRepository;
use Doctrine\ORM\EntityManagerInterface;
use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputArgument as InputArgumentAlias;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\PasswordHasher\Hasher\UserPasswordHasherInterface;

#[AsCommand(
    name: 'register:user',
    description: 'Register a new user',
    hidden: false,
)]
class SeedUsers extends Command
{
    private UserRepository $userRepository;
    private UserPasswordHasherInterface $passwordHasher;
    private EntityManagerInterface $entityManager;

    public function __construct(
        EntityManagerInterface $entityManager,
        UserRepository $userRepository,
        UserPasswordHasherInterface $passwordHasher)
    {
        parent::__construct();
        $this->entityManager = $entityManager;
        $this->userRepository = $userRepository;
        $this->passwordHasher = $passwordHasher;
    }

    public function configure(): void
    {
        $this
            ->addArgument('email', InputArgumentAlias::REQUIRED, "User's email")
            ->addArgument('password', InputArgumentAlias::REQUIRED, "User's password");
    }

    public function execute(InputInterface $input, OutputInterface $output): int
    {
        $output->writeln('<info>Seeding users</info>');
        $email = $input->getArgument('email');
        $password = $input->getArgument('password');
        $output->writeln('Email: ' . $email);
        $output->writeln('Password: ' . $password);

        $user = new User();
        $user->setEmail($email);
        $user->setPassword($this->passwordHasher->hashPassword($user, $password));
        $this->entityManager->persist($user);
        $this->entityManager->flush();

        return Command::SUCCESS;
    }
}
