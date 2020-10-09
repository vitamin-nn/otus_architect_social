<template>
  <div class="container">
      <h3>Main page</h3>
      <div class="card mb-3" v-for="item in items" v-bind:key="item.id">
        <div class="card-body" v-if="item.id">
          <h5 class="card-title">{{ item.first_name }} {{ item.last_name }}</h5>
          <p class="card-text">Email: {{ item.email }}</p>
          <a v-bind:href="'/user/' + item.id"  class="btn btn-primary">Go to profile</a>
        </div>
      </div>

  </div>
</template>

<script>
import UserService from "../services/user.service";

export default {
  name: "Home",
  data: () => ({
    items: [],
  }),

  created() {
    UserService.getMainPageContent().then(
      (response) => {
        this.items = response.data;
      },
      (error) => {
        this.content =
          (error.response && error.response.data) ||
          error.message ||
          error.error.toString();
      }
    );
  },
};
</script>