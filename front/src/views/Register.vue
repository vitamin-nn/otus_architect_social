<template>
  <div class="col-md-12">
    <div class="card card-container">
      <img
        id="profile-img"
        src="//ssl.gstatic.com/accounts/ui/avatar_2x.png"
        class="profile-img-card"
      />
      <form name="form" @submit.prevent="handleRegister">
        <div v-if="!successful">
          <div class="form-group">
            <label for="email">Email</label>
            <input
              v-model="user.email"
              v-validate="'required|email|max:50'"
              type="email"
              class="form-control"
              name="email"
            />
            <div v-if="submitted && errors.has('email')" class="alert-danger">
              {{ errors.first("email") }}
            </div>
          </div>
          <div class="form-group">
            <label for="password">Password</label>
            <input
              v-model="user.password"
              v-validate="'required|min:6|max:40'"
              type="password"
              class="form-control"
              name="password"
            />
            <div
              v-if="submitted && errors.has('password')"
              class="alert-danger"
            >
              {{ errors.first("password") }}
            </div>
          </div>
          <div class="form-group">
            <label for="first_name">First name</label>
            <input
              v-model="user.first_name"
              v-validate="'required|min:2|max:50'"
              type="text"
              class="form-control"
              name="first_name"
            />
            <div
              v-if="submitted && errors.has('first_name')"
              class="alert-danger"
            >
              {{ errors.first("first_name") }}
            </div>
          </div>
          <div class="form-group">
            <label for="last_name">Last name</label>
            <input
              v-model="user.last_name"
              v-validate="'required|min:2|max:50'"
              type="text"
              class="form-control"
              name="last_name"
            />
            <div
              v-if="submitted && errors.has('last_name')"
              class="alert-danger"
            >
              {{ errors.first("last_name") }}
            </div>
          </div>
          <div class="form-group">
            <label for="birth_date">Birth date (YYYY-MM-DDT00:00:00Z)</label>
            <input
              v-model="user.birth_date"
              v-validate="'required|min:20|max:20'"
              type="text"
              class="form-control"
              name="birth_date"
            />
            <div
              v-if="submitted && errors.has('birth_date')"
              class="alert-danger"
            >
              {{ errors.first("birth_date") }}
            </div>
          </div>
          <div class="form-group">
            <label for="sex">Sex (M|F)</label>
            <input
              v-model="user.sex"
              v-validate="'required|min:1|max:1'"
              type="text"
              class="form-control"
              name="sex"
            />
            <div v-if="submitted && errors.has('sex')" class="alert-danger">
              {{ errors.first("sex") }}
            </div>
          </div>
          <div class="form-group">
            <label for="interest">Interests</label>
            <input
              v-model="user.interest"
              v-validate="'required|min:0|max:1024'"
              type="text"
              class="form-control"
              name="interest"
            />
            <div
              v-if="submitted && errors.has('interest')"
              class="alert-danger"
            >
              {{ errors.first("interest") }}
            </div>
          </div>
          <div class="form-group">
            <label for="city">City</label>
            <input
              v-model="user.city"
              v-validate="'required|min:0|max:100'"
              type="text"
              class="form-control"
              name="city"
            />
            <div v-if="submitted && errors.has('city')" class="alert-danger">
              {{ errors.first("city") }}
            </div>
          </div>
          <div class="form-group">
            <button class="btn btn-primary btn-block">Sign Up</button>
          </div>
        </div>
      </form>

      <div
        v-if="message"
        class="alert"
        :class="successful ? 'alert-success' : 'alert-danger'"
      >
        {{ message }}
      </div>
    </div>
  </div>
</template>

<script>
import User from "../models/user";

export default {
  name: "Register",
  data() {
    return {
      user: new User("", "", ""),
      submitted: false,
      successful: false,
      message: "",
    };
  },
  computed: {
    loggedIn() {
      return this.$store.state.auth.status.loggedIn;
    },
  },
  mounted() {
    if (this.loggedIn) {
      this.$router.push("/profile");
    }
  },
  created() {
    if (this.loggedIn) {
      this.$router.push("/profile");
    }
  },
  methods: {
    handleRegister() {
      this.message = "";
      this.submitted = true;
      this.$validator.validate().then((isValid) => {
        if (isValid) {
          this.$store.dispatch("auth/register", this.user).then(
            () => {
              this.$router.push("/profile");
            },
            (error) => {
              this.message =
                (error.response && error.response.data) ||
                error.message ||
                error.error.toString();
              this.successful = false;
            }
          );
        }
      });
    },
  },
};
</script>

<style scoped>
label {
  display: block;
  margin-top: 10px;
}

.card-container.card {
  max-width: 350px !important;
  padding: 40px 40px;
}

.card {
  background-color: #f7f7f7;
  padding: 20px 25px 30px;
  margin: 0 auto 25px;
  margin-top: 50px;
  -moz-border-radius: 2px;
  -webkit-border-radius: 2px;
  border-radius: 2px;
  -moz-box-shadow: 0px 2px 2px rgba(0, 0, 0, 0.3);
  -webkit-box-shadow: 0px 2px 2px rgba(0, 0, 0, 0.3);
  box-shadow: 0px 2px 2px rgba(0, 0, 0, 0.3);
}

.profile-img-card {
  width: 96px;
  height: 96px;
  margin: 0 auto 10px;
  display: block;
  -moz-border-radius: 50%;
  -webkit-border-radius: 50%;
  border-radius: 50%;
}
</style>