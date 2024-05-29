<script setup>

</script>

<template>
  <div class="body">
    <div class="header">

    </div>
    <div class="main">
      <div class="side-bar">

      </div>
      <div class="main-content">
        <div class="card w-96 glass">
            <figure><img alt="" src="https://img.daisyui.com/images/stock/photo-1534528741775-53994a69daeb.jpg" /></figure>
          <div class="card-body">
            <h2 class="card-title">Name Surname</h2>
            <p>Role: client</p>
            <div class="stats shadow">

              <div class="stat">
                <div class="stat-title">Total of Sport Activities</div>
                <div class="stat-value">124 days</div>

              </div>

            </div>

          </div>
        </div>
        <div class="info">
          <h2 class="card-title">Информация</h2>
          <label class="input input-bordered flex items-center gap-2">
          Name
            <input type="text" class="grow"  />
        </label>
          <label class="input input-bordered flex items-center gap-2">
            Surname
            <input type="text" class="grow"  />
          </label>
          <label class="input input-bordered flex items-center gap-2">
            Email
            <input type="text" class="grow"  />
          </label>



        </div>




      </div>
    </div>

  </div>
</template>

<style scoped>
.body{
  background-color:#16141C;
  width: 100vw;
  height: 100vh;


}
.header
{
  background-color:#6E52F9;
  width:100%;
  height: 10vh;
  text-align: center;
}

.main
{
  width:100%;
  height:90vh;
  display: flex;
}
.side-bar{
  background-color: #6E52F9;
  margin: 1.5% 0;
  width: 12%;
  border-bottom-right-radius: 30px;
  border-top-right-radius: 30px;

}

.main-content
{
  background-color: #6E52F9;
  margin: 1.5% 2% 1.5% 3%;
  width: 83%;
  display: flex;
}
.card{
  margin: 3% 3%;
}
.info{
  margin: 3% 3%;
  display: block;
}
.info >.input{
  margin-top: 1%;
}
.calendar{


}


</style>
