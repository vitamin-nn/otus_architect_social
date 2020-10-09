<template>
  <div class="container">
    <h3>Profile</h3>
    <ul class="list-group list-group-flush">
      <li class="list-group-item"><strong>Email:</strong> {{ profile.email}}</li>
      <li class="list-group-item"><strong>First name:</strong> {{ profile.first_name}}</li>
      <li class="list-group-item"><strong>Last name:</strong> {{ profile.last_name}}</li>
      <li class="list-group-item"><strong>Sex:</strong> {{ profile.sex}}</li>
      <li class="list-group-item"><strong>Interest:</strong> {{ profile.interest}}</li>
      <li class="list-group-item"><strong>City:</strong> {{ profile.city}}</li>
    </ul>
    <div class="form-group">
      <div v-if="message" class="alert alert-danger" role="alert">
        {{ message }}
      </div>
      <div v-if="success" class="alert alert-success" role="alert">
        Success!
      </div>
    </div>
    <template v-if="currentUser && currentUser.id != profile.id">
      <button class="btn btn-primary btn-block" :disabled="sent" v-on:click="addToFriends">
        <span v-show="loading" class="spinner-border spinner-border-sm"></span>
        Add to friends
      </button>
    </template>
  </div>
</template>

<script>
import UserService from "../services/user.service";

export default {
  name: "User",
  data: () => ({
    profile: [],
    loading: false,
    sent: false,
    success: false,
    message: "",
  }),

  computed: {
    currentUser() {
      return this.$store.state.auth.user;
    },
  },
  created() {
    this.loading = true;
    UserService.getPublicProfileById(this.$route.params.userId).then(
      (response) => {
        this.profile = response.data;
        this.loading = false;
      },
      (error) => {
        this.message =
          (error.response && error.response.data) ||
          error.message ||
          error.error.toString();
          this.loading = false;
      }
    );
  },
  methods: {
    addToFriends: function () {
      this.loading = true;
      UserService.addFriend(this.profile.id).then(
        () => {
          this.loading = false;
          this.success = true;
          this.sent = true;
        },
        (error) => {
          this.message =
            (error.response && error.response.data.error);
            this.loading = false;
        }
      );
    }
  }
};
</script>