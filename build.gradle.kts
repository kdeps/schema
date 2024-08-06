plugins {
  id("org.pkl-lang") version "0.25.2"
}

val maybeVersion = System.getenv("VERSION")

pkl {
  project {
    packagers {
      register("makePackages") {
        if (maybeVersion != null) {
          environmentVariables.put("VERSION", maybeVersion)
        }
        projectDirectories.from(file("schema/"))
      }
    }
  }
  // ./gradlew pkldoc
  if (maybeVersion != null) {
    pkldocGenerators {
      register("pkldoc") {
        sourceModules =
          listOf(uri("package://schema.kdeps.com/core@$maybeVersion"))
      }
    }
  }
}
