<script setup>
import {  ref  } from 'vue';
import { useRouter } from 'vue-router';
import LoginIcon from '@/components/icons/LoginIcons.vue'
import PassIcons from '@/components/icons/PassIcons.vue'
import ShowIcons from '@/components/icons/ShowIcons.vue'
import HideIcons from '@/components/icons/HideIcons.vue'


const password = ref('');
const showPassword = ref(false);
const router = useRouter();

const togglePasswordVisibility = () =>{
    showPassword.value = !showPassword.value;
};
const navigateToReg = ()=>{
    router.push('/reg')
};

</script>


<template>
    <div class="main">
        
        
            <div class="container-auth">
               
                <div class="container">
                    <div class="logo">
                        <p>Авторизация</p>
                    </div>
                    
                    <div class="field-input">
                        <span>Username</span>
                        <LoginIcon/>
                        <input class="input" type="email" placeholder="Type your username">
                    </div>
                    
                    <div class="field-input">
                        <span>Password</span>
                        <PassIcons/>
                        <input class="input"  v-model="password" :type="showPassword ? 'text' : 'password'"  placeholder="Type your password"  >
                        <button class="showHide" @click="togglePasswordVisibility">
                        <HideIcons v-if="showPassword" />
                        <ShowIcons v-else />
                       
                        </button>
                    </div>
                    
                    <div class="forgot-pass">
                        <a href="#">Forgot password?</a>
                    </div>
                    <div class="remember">
                        <input  class="input-checkbox" type="checkbox" id="remember" name="Remember me" placeholder="Remember me">
                        <label for="remember">Remember me</label>
                    </div>
                    <div class="link-reg">
                        <p>Dont have an account? </p>
                        <a @click="navigateToReg">Click here</a>
                    </div>
                    <div class="button">
                        <button class="button-auth" type="submit">Войти</button>
                    </div>
                  
                    

                </div>

               
            </div>
        

    </div>
</template>



<style scoped>

   *{
    font-family: Inter;
   }
    .main{
       background-color:#16141C; 
       display: flex;
       align-items: center;
       justify-content: center;
       width: 100vw;
       height: 100vh;
    }
      
    .container-auth
    {
        height: 45vh;
        width: 17vw;
        background-color: #24212B;
        border-radius: 10%;
        display: block;
        position: relative;
        overflow: hidden;
        
    }
    .container{
        margin: 10% 10%;
        padding: 0 auto;
       
        
    }
    .logo{
        width: 100%;
        align-items: center;
        overflow: hidden;
    }
    p{
        color:white;
        text-align: center;
        font-size: 1.5em;
        margin-bottom:  2% ;


    }
    .field-input{
        width: 100%;
        height: 100%;
        border-bottom: 2px solid #FFF;
        position: relative;
        align-items: center;
        margin-top: 10px;
        
       

    }
    .field-input span{
        font-size: 1em ;
        color: #FFF;
        padding: 0.3em auto 20px auto;
    }
    .input{
        display: block;
        width:  70%;
        height: 100%;
        background-color: #24212B;
        margin: 10px 25px;
       
       

    }
    .showHide
    {
        outline: none;
        border: none;
        background: transparent;
        position: absolute;
        position: absolute;
        top:65%;
        right:0px;
        transform: translateY(-50%);
    }
    .forgot-pass{
        width: 100%;
        margin:15px auto
        
    }
    .forgot-pass a{
        text-align: right
    }
    
    .container input{
        outline:nine;
        border: none;
        font-size:1em;
        color: white;

       } 
       input::placeholder
       {
        justify-content: center;
        color: white;
        font-size: 1em;
        vertical-align:middle
        
       }
       input:focus{
        color: white;
       
}
       a{
        display: block;
        text-align: left;
        font-size: 1em;
        text-decoration: none;
        color: #fff;
        
       }
       a:hover
       {
        color: #6E52F9;
       }
       .input-checkbox{
        margin-right: 10px;
        transform: scale(1.5);
        
       }
       .input-checkbox+label
       {
        display: inline-block;
        font-size: 1em;
        color:#FFF;
        cursor:pointer;
       }
       .link-reg{
        margin-top: 5%;
        display: flex;
       }
       .link-reg p{
        font-size: 1em;
       }
       .button{
        text-align: center;
        margin: 10% auto;
        width: 100%;
    
       }
       .button-auth{
        color:#FFFFFF;
        justify-content: center;
        text-align: center;
        background-color: #6E52F9;
        border:none;
        border-radius: 1.25em;
        padding: 0.75em 2em;
        font-size: 1em;
        min-width: 35%;
        


       }
       


</style>