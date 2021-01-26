<template>
  <div class="container">
    <h3>My Profile</h3>
    <ul class="list-group list-group-flush">
      <li class="list-group-item"><strong>Id:</strong> {{ profile.id}}</li>
      <li class="list-group-item"><strong>Email:</strong> {{ profile.email}}</li>
      <li class="list-group-item"><strong>First name:</strong> {{ profile.first_name}}</li>
      <li class="list-group-item"><strong>Last name:</strong> {{ profile.last_name}}</li>
      <li v-if="profile.birth_date" class="list-group-item"><strong>Birth date:</strong> {{ displayDate(profile.birth_date) }}</li>
      <li v-if="profile.sex" class="list-group-item"><strong>Sex:</strong> {{ profile.sex}}</li>
      <li v-if="profile.interest" class="list-group-item"><strong>Interest:</strong> {{ profile.interest}}</li>
      <li v-if="profile.city" class="list-group-item"><strong>City:</strong> {{ profile.city}}</li>
    </ul>
    <form name="formFeed" @submit.prevent="handleFeed">
      <div class="input-group">
        <span class="input-group-text">Add feed message</span>
        <textarea class="form-control" aria-label="With textarea" v-model="feedInputMessage"></textarea>
      </div>
      <div class="form-group">
        <button class="btn btn-primary btn-block" :disabled="loading">
          <span
            v-show="loading"
            class="spinner-border spinner-border-sm"
          ></span>
          <span>Save</span>
        </button>
      </div>
    </form>
    <div class="card" v-for="item in feedItems" v-bind:key="item.id">
      <div class="card-body">
        {{ item.body }} (from user id: {{item.user_id}})
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
import FeedService from "../services/feed.service";

export default {
  name: "Profile",
  data: () => ({
    profile: [],
    feedItems: [],
    message: "",
    loading: false,
  }),

  computed: {
    currentUser() {
      return this.$store.state.auth.user;
    },
  },
  mounted() {
    if (!this.currentUser) {
      this.$router.push("/login");
    }
  },
  created() {
    UserService.getMyProfile().then(
      (response) => {
        this.profile = response.data;
      },
      (error) => {
        this.message =
          (error.response && error.response.data) ||
          error.message ||
          error.error.toString();
      }
    );

    FeedService.getFeed(this.currentUser.id).then(
      (response) => {
        this.feedItems = response.data;
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
    displayDate: function(date) {
      var d = new Date(date)
      return d.getDate() + '.' + (d.getMonth() + 1) + '.' + d.getFullYear()
    },
    handleFeed() {
      this.loading = true;
      this.$validator.validateAll().then((isValid) => {
        if (!isValid) {
          this.loading = false;
          return;
        }

        if (this.feedInputMessage) {
          FeedService.addFeed(this.feedInputMessage).then(
            () => {
              window.location.reload();
            },
            (error) => {
              this.loading = false;
              this.message =
                (error.response && error.response.data) ||
                error.message ||
                error.error.toString();
            }
          );
        }
      });
    },
  }
};
</script>