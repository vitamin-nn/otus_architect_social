<template>
  <div class="container">
    <h3>My friends</h3>
    <div class="card mb-3" v-for="item in items" v-bind:key="item.id">
        <div class="card-body">
            <h5 class="card-title">{{ item.first_name }} {{ item.last_name }}</h5>
            <p class="card-text">Email: {{ item.email }}</p>
            <template v-if="currentUser.id != item.id">
                <button class="btn btn-primary btn-block" :disabled="loading" v-on:click="removeFriend(item.id)">
                    <span v-show="loading" class="spinner-border spinner-border-sm"></span>
                    Remove from friends
                </button>
            </template>
        </div>
    </div>
    <div class="form-group">
        <div v-if="message" class="alert alert-danger" role="alert">
            {{ message }}
        </div>
    </div>
  </div>
</template>

<script>
import UserService from "../services/user.service";

export default {
  name: "Friends",
  data: () => ({
    items: [],
    loading: false,
    message: "",
  }),
  computed: {
    currentUser() {
      return this.$store.state.auth.user;
    },
  },
  created() {
    UserService.getFriendList().then(
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
  methods: {
    removeFriend: function (id) {
      this.loading = true;
      UserService.removeFriend(id).then(
        () => {
          this.loading = false;
          location.reload()
          //this.profile = response.data;
        },
        (error) => {
          this.message =
            (error.response && error.response.data);
            this.loading = false;
        }
      );
    }
  }
};
</script>
