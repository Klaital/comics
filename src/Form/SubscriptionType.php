<?php

namespace App\Form;

use App\Entity\Subscription;
use Symfony\Component\Form\AbstractType;
use Symfony\Component\Form\Extension\Core\Type\CheckboxType;
use Symfony\Component\Form\Extension\Core\Type\IntegerType;
use Symfony\Component\Form\Extension\Core\Type\SubmitType;
use Symfony\Component\Form\Extension\Core\Type\TextType;
use Symfony\Component\Form\FormBuilderInterface;
use Symfony\Component\OptionsResolver\OptionsResolver;

class SubscriptionType extends AbstractType
{
    public function buildForm(FormBuilderInterface $builder, array $options): void
    {
        $builder
            ->add('title', TextType::class, [
                'label' => 'Title',
            ])
            ->add('homepage', TextType::class, [
                'label' => 'Homepage',
            ])
            ->add('rssUrl', TextType::class,  [
                'required' => false,
                'label' => 'RSS URL',
            ])
            ->add('ordinal', IntegerType::class)
            ->add('firstComicUrl', TextType::class, [
                'required' => false,
                'label' => 'First Comic URL',
            ])
            ->add('lastestComicUrl', TextType::class, [
                'required' => false,
                'label' => 'Latest Comic URL',
            ])
            ->add('updatesMon', CheckboxType::class, [
                'required' => false,
                'label' => 'Updates Monday?',
            ])
            ->add('updatesTue', CheckboxType::class, [
                'required' => false,
                'label' => 'Updates Tuesday?',
            ])
            ->add('updatesWed', CheckboxType::class, [
                'required' => false,
                'label' => 'Updates Wednesday?',
            ])
            ->add('updatesThu', CheckboxType::class, [
                'required' => false,
                'label' => 'Updates Thursday?',
            ])
            ->add('updatesFri', CheckboxType::class, [
                'required' => false,
                'label' => 'Updates Friday?',
            ])
            ->add('updatesSat', CheckboxType::class, [
                'required' => false,
                'label' => 'Updates Saturday?',
            ])
            ->add('updatesSun', CheckboxType::class, [
                'required' => false,
                'label' => 'Updates Sunday?',
            ])
            ->add('active', CheckboxType::class, [
                'required' => false,
                'label' => 'Is Active?',
            ])
            ->add('nsfw', CheckboxType::class, [
                'required' => false,
                'label' => 'Is NSFW?',
            ])
            ->add('submit', SubmitType::class);
    }

    public function configureOptions(OptionsResolver $resolver): void
    {
        $resolver->setDefaults([
            'data_class' => Subscription::class,
        ]);
    }
}
