<template>
  <div class="security-settings">
    <h2>Security Settings</h2>
    <div class="form-group">
      <label for="current-password">Current Password</label>
      <input type="password" id="current-password" placeholder="Current Password" />
    </div>
    <div class="form-group">
      <label for="new-password">New Password</label>
      <input type="password" id="new-password" placeholder="New Password" />
    </div>
    <div class="form-group">
      <label for="confirm-password">Confirm New Password</label>
      <input type="password" id="confirm-password" placeholder="Confirm New Password" />
    </div>
    <div class="two-factor-auth">
      <h3><i class="fas fa-shield-alt"></i> Two-factor authentication</h3>
      <div class="two-factor-content">
        <img src="" alt="Two-factor authentication" />
        <div class="two-factor-description">
          <p>Two factor authentication is not enabled yet.</p>
          <p>
            Two factor authentication adds an additional layer of security to your account by
            requiring more than just a password to log in. Learn more.
          </p>
        </div>
      </div>
    </div>
    <div class="actions">
      <button class="save-button">Save Changes</button>
      <button class="cancel-button">Cancel</button>
    </div>
  </div>
</template>

<script setup></script>

<style scoped>
.security-settings {
  background-color: #2e3440;
  padding: 20px;
  border-radius: 10px;
  color: #eceff4;
}

h2 {
  color: #88c0d0;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  color: #eceff4;
}

.form-group input {
  width: 100%;
  padding: 10px;
  border: 1px solid #4c566a;
  border-radius: 5px;
  background-color: #434c5e;
  color: #eceff4;
}

.two-factor-auth {
  margin-top: 20px;
}

.two-factor-auth h3 {
  display: flex;
  align-items: center;
  color: #88c0d0;
}

.two-factor-auth i {
  margin-right: 10px;
}

.two-factor-content {
  display: flex;
  align-items: center;
  margin-top: 10px;
}

.two-factor-content img {
  width: 150px;
  height: auto;
  margin-right: 20px;
}

.two-factor-description {
  color: #d8dee9;
}

.actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.save-button,
.cancel-button {
  padding: 10px 20px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.save-button {
  background-color: #5e81ac;
  color: #eceff4;
  margin-right: 10px;
}

.cancel-button {
  background-color: #d08770;
  color: #eceff4;
}
</style>
