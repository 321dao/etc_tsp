<template>
  <div class="login_form" :style="background">
    <Form class="inner_form" ref="formInline" :model="formInline" :rules="ruleInline">
      <FormItem prop="user">
        <Input
          prefix="ios-person-outline"
          type="text"
          v-model="formInline.user"
          placeholder="Username"
        />
      </FormItem>
      <FormItem prop="password">
        <Input
          prefix="ios-lock-outline"
          type="password"
          v-model="formInline.password"
          placeholder="Password"
        />
      </FormItem>
      <FormItem>
        <Button type="primary" @click="handleSubmit('formInline')">Login</Button>
      </FormItem>
    </Form>
  </div>
</template>
<script>
export default {
  data() {
    return {
      formInline: {
        user: "",
        password: ""
      },
      ruleInline: {
        user: [
          {
            required: true,
            message: "Please fill in the user name",
            trigger: "blur"
          }
        ],
        password: [
          {
            required: true,
            message: "Please fill in the password.",
            trigger: "blur"
          },
          {
            type: "string",
            min: 6,
            message: "The password length cannot be less than 6 bits",
            trigger: "blur"
          }
        ]
      },
      background: {
        backgroundImage: "url(" + require("../assets/login_bg.jpg") + ")",
        backgroundRepeat: "no-repeat",
        backgroundSize: "cover",
        marginTop: "0px"
      }
    };
  },
  methods: {
    handleSubmit(name) {
      this.$refs[name].validate(valid => {
        if (valid) {
          this.axios
            .post("/api/login", {
              user: this.formInline.user,
              password: this.formInline.password
            })
            .then(response => {
              this.$Message.success("Success!");
              window.console.log(response.data);
              window.localStorage["token"] = response.data.token;

              this.$router.push("/mainpage/devices");
            })
            .catch(error => {
              this.$Message.error("usrname or password is error!");
              window.console.log(error);
            });
        } else {
          this.$Message.error("Fail!");
        }
      });
    }
  }
};
</script>

<style scoped>
.login_form {
  display: flex;
  justify-content: center; /* 水平居中 */
  align-items: center; /* 垂直居中 */
  width: 100%;
  height: 100%;
  position: absolute;
}
.inner_form {
  width: 300px;
  height: 250px;
}
</style>