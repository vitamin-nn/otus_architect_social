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
  </div>
</template>

<script>
import UserService from "../services/user.service";

export default {
  name: "Profile",
  data: () => ({
    profile: [],
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
        this.content =
          (error.response && error.response.data) ||
          error.message ||
          error.error.toString();
      }
    );
  },
  methods: {
    displayDate: function(date) {
      var d = new Date(date)
      return d.getDate() + '.' + (d.getMonth() + 1) + '.' + d.getFullYear()
    }
  }
};
</script>