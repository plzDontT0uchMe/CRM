export const subscriptions = [
  {
    id: 1,
    name: 'basic',
    price: {
      oneMonth: 'Free',
      threeMonths: 'Free',
      sixMonths: 'Free',
      twelveMonths: 'Free'
    },
    description:
      'Use your own workout program, maintain a nutrition diary, and maintain a workout diary.',
    possibilities: ['workoutProgram', 'nutritionDiary', 'workoutDiary'],
    color: '#ffcccb'
  },
  {
    id: 2,
    name: 'standard',
    price: {
      oneMonth: '$40/month',
      threeMonths: '$100/3 months',
      sixMonths: '$180/6 months',
      twelveMonths: '$320/year'
    },
    description:
      'Includes all basic features plus the ability to book gym time (specify time and date).',
    possibilities: ['workoutProgram', 'nutritionDiary', 'workoutDiary', 'bookGym'],
    color: '#add8e6'
  },
  {
    id: 3,
    name: 'premium',
    price: {
      oneMonth: '$60/month',
      threeMonths: '$150/3 months',
      sixMonths: '$270/6 months',
      twelveMonths: '$480/year'
    },
    description:
      'Includes all standard features plus the ability to book a trainer (choose your trainer).',
    possibilities: ['workoutProgram', 'nutritionDiary', 'workoutDiary', 'bookGym', 'bookTrainer'],
    color: '#90ee90'
  }
]
