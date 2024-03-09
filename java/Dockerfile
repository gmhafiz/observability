FROM gradle:jdk17 as build
COPY --chown=gradle:gradle . /home/gradle/src
WORKDIR /home/gradle/src
RUN gradle build --no-daemon
#RUN gradle build


# https://spring.io/guides/topicals/spring-boot-docker/
FROM openjdk:17

#VOLUME /tmp
RUN mkdir /app

#ARG JAR_FILE
#COPY ${JAR_FILE} /tmp/app.jar

COPY --from=build /home/gradle/src/build/libs/demo-0.0.1-SNAPSHOT.jar app.jar
#COPY --from=build /home/gradle/src/build/libs/*.jar app.jar
COPY opentelemetry-javaagent.jar .

EXPOSE 8080

ENTRYPOINT ["java", "-javaagent:opentelemetry-javaagent.jar", "-jar", "app.jar"]