<template>
<div class="main">  
    <div class="main-content">
      <p>Управление ролями пользователя</p>
      <div class="search">
        <input type="text" v-model="searchQuery" placeholder="Введите логин для поиска"/>
        <!-- <button class="search-btn" @click="search">
           <SearchIcons /> 
         </button> -->
      </div>
      
      
      <div class="card" v-for="user in filtredUsers" :key="user.id">
        <p>{{ user.name }} {{ user.surname }} {{ user.patronomic }} Логин:{{ user.login }} - Роль: {{ user.role }}</p>
        <select class="select" v-model="user.role" @change="updateRole(user.id, $event.target.value)">
          <option value="Администратор">Администратор</option>
          <option value="Тренер">Тренер</option>
          <option value="Клиент">Клиент</option>
        </select>
      </div>
    </div>
</div>
  </template>
  
  <script setup>
  import { computed, ref } from 'vue';
  import SearchIcons from '@/components/icons/SearchIcons.vue'
  
 
      const users = ref([
        { id: 1, login: 'qwer', name: 'User1', surname:'Surname1', patronomic:'Patronomic1', role: 'Клиент' },
        { id: 2, login: 'asd', name: 'User2', surname:'Surname2', patronomic:'Patronomic2', role: 'Клиент' },
        { id: 3, login: 'qwd', name: 'User3', surname:'Surname3', patronomic:'Patronomic3', role: 'Тренер' },
        { id: 4, login: '1234', name: 'User4', surname:'Surname4', patronomic:'Patronomic4', role: 'Администратор' },
       
      ]);
    const searchQuery = ref('');

      const filtredUsers = computed(()=>{
        if(!searchQuery.value){
            return users.value;
        }else
        return users.value.filter(user=>{
            return user.login.toLowerCase().includes(searchQuery.value.toLowerCase());
        });
      });
      const updateRole = (userId, newRole) => {
        const user = users.value.find(user => user.id === userId);
        if (user) {
          user.role = newRole;
          
        }
      };
  
  </script>
  <style  scoped>
  .main{
    background-color:#16141C; 
       display: flex;
       align-items: center;
       justify-content: center;
       width: 100vw;
       height: 100vh;
  }

  .main-content{
        height: 85vh;
        width: 87vw;
        background-color: #24212B;
        border-radius: 15px;
        display: block;
        overflow: hidden;
        align-items: center;
        text-align: center;
        
  }
  p{
    font-size: 20px;
    padding: 15px;
  }
  .search{
    width: 100%;
  }
  .search input{
    outline: none;
    border: none;
    border-radius: 10px;
    font-size:1em;
    color: white;
    width: 80%;
    height: 35px;
    background-color: #6E52F9;

  }
  .search input::placeholder
       {
        justify-content: center;
        color: white;
        font-size: 1em;
        vertical-align:middle
       }
  .search-btn{
        outline: none;
        border: none;
        background: transparent;
        position: absolute;
        position: absolute;
        top:12%;
        right:16%;
        transform: translateY(-50%);

  }
  p{
    font-size: 20px;
    color: white;

  }
  .card{
    width: 90%;
    height: 10%;
    background-color:#6E52F9;
    border-radius: 30px;
    margin: 20px;
    display: flex;
    align-items: center;
    padding:0 20px;
    position:relative ;
  }
  .card p{
    font-size: 20px;
    color: white;
  }
  .select{
    border: none;
    border-radius:10px;
    font-size: 18px;
    color: white;
    background-color: #24212B;
    position: absolute;
    right: 10px;

  }



  
  </style>