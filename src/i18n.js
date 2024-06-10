import { createI18n } from 'vue-i18n'

const messages = {
  en: {
    calendar: 'Calendar',
    nutritionPlans: 'Nutrition Plans',
    workoutPlans: 'Workout Plans',
    createNewPlan: 'Create new plan',
    newPlan: 'New plan',
    planName: 'Plan name',
    planDescription: 'Plan description',
    planExercises: 'Plan exercises',
    addExercise: 'Add exercise',
    selectedExercises: 'Selected exercises',
    addPlan: 'Add plan',
    cancel: 'Cancel',
    loading: 'Loading',
    chat: 'Chat',
    exercises: 'Exercises',
    profile: 'Profile',
    settings: 'Settings',
    account: 'Account',
    security: 'Security',
    info: 'Info',
    notifications: 'Notifications',
    logout: 'Logout',
    title: 'Title',
    titlePlaceholder: 'Enter the title of the event',
    start: 'Start',
    startPlaceholder: 'Select the start date and time',
    end: 'End',
    endPlaceholder: 'Select the end date and time',
    save: 'Save',
    cancel: 'Cancel',
    memberShip: 'Membership',
    chooseTrainer: 'Choose Trainer',
    age: 'Age',
    experience: 'Experience',
    gender: 'Gender',
    sportType: 'Sport type',
    choose: 'Choose',
    achievements: 'Achievements',
    buy: 'Buy',
    firstName: 'Name',
    lastName: 'Surname',
    middleName: 'Patronomic',
    dob: 'Date of Birth',
    accountCreated: 'Account created',
    position: 'Role',

    durations: {
      oneMonth: '1 month',
      threeMonths: '3 months',
      sixMonths: '6 months',
      twelveMonths: '12 months'
    },
    features: {
      workoutProgram: 'Use your own workout program',
      nutritionDiary: 'Maintain a nutrition diary',
      workoutDiary: 'Maintain a workout diary',
      bookGym: 'Book gym time (specify time and date)',
      bookTrainer: 'Book a trainer (choose your trainer)'
    },
    plans: {
      basic: 'Basic Plan',
      standard: 'Standard Plan',
      premium: 'Premium Plan'
    }
  },
  ru: {
    calendar: 'Календарь',
    nutritionPlans: 'Планы питания',
    workoutPlans: 'Планы тренировок',
    createNewPlan: 'Создать новый план',
    newPlan: 'Новый план',
    planName: 'Название плана',
    planDescription: 'Описание плана',
    planExercises: 'План упражнений',
    addExercise: 'Добавить упражнения',
    selectedExercises: 'Выбранные кпражнения',
    addPlan: 'Добавить план',
    cancel: 'Назад',
    loading: 'Загрузка',
    chat: 'Чат',
    exercises: 'Упражнения',
    profile: 'Профиль',
    account: 'Аккаунт',
    security: 'Защита',
    info: 'Информация',
    notifications: 'Уведомления',
    settings: 'Настройки',
    logout: 'Выход',
    title: 'Название',
    titlePlaceholder: 'Введите название события',
    start: 'Начало',
    startPlaceholder: 'Выберите дату и время начала',
    end: 'Конец',
    endPlaceholder: 'Выберите дату и время окончания',
    save: 'Сохранить',
    cancel: 'Отмена',
    memberShip: 'Абонементы',
    chooseTrainer: 'Выбрать тренера',
    age: 'Возраст',
    experience: 'Стаж работы',
    gender: 'Пол',
    sportType: 'Вид спорта',
    achievements: 'Достижения',
    choose: 'Выбрать',
    buy: 'Купить',
    firstName: 'Имя',
    lastName: 'Фамилия',
    middleName: 'Отчество',
    dob: 'Дата рождения',
    accountCreated: 'Аккаунт создан',
    position: 'Роль',
    durations: {
      oneMonth: '1 месяц',
      threeMonths: '3 месяца',
      sixMonths: '6 месяцев',
      twelveMonths: '12 месяцев'
    },
    features: {
      workoutProgram: 'Использовать собственную программу тренировок',
      nutritionDiary: 'Вести дневник питания',
      workoutDiary: 'Вести дневник тренировок',
      bookGym: 'Забронировать время в зале (указать время и дату)',
      bookTrainer: 'Забронировать тренера (выбрать тренера)'
    },
    plans: {
      basic: 'Базовый план',
      standard: 'Стандартный план',
      premium: 'Премиум план'
    }
  }
}

const i18n = createI18n({
  locale: 'en',
  fallbackLocale: 'en',
  messages
})

export default i18n
