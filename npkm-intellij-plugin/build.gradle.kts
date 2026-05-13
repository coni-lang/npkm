plugins {
    id("java")
    id("org.jetbrains.intellij") version "1.16.0"
}

group = "com.hellonico.npkm"
version = "1.0.0"

repositories {
    mavenCentral()
}

intellij {
    version.set("2023.2.5")
    type.set("IC")
    plugins.set(listOf("com.intellij.java", "yaml"))
}

tasks {
    withType<JavaCompile> {
        sourceCompatibility = "17"
        targetCompatibility = "17"
    }
}
