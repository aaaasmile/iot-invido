<template>
  <v-row justify="center">
    <v-col class="mb-12" cols="12" md="10">
      <v-card>
        <v-card-title>Login</v-card-title>
        <v-container>
          <v-row justify="space-around">
            <v-col cols="8">
              <v-row justify="space-around">
                <v-text-field
                  label="Username"
                  v-model="Username"
                ></v-text-field>
              </v-row>
              <v-row>
                <v-text-field
                  label="Password"
                  type="password"
                  v-model="password"
                ></v-text-field>
              </v-row>
            </v-col>
          </v-row>
        </v-container>
        <v-card-actions>
          <v-row justify="space-around">
            <v-btn color="primary" text @click="signIn"> Sign In </v-btn>
          </v-row>
        </v-card-actions>
      </v-card>

      <v-row v-if="IsLastMessage">
        <v-container>
          {{ LastMsg }}
        </v-container>
      </v-row>
    </v-col>
  </v-row>
</template>