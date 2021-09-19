FROM scratch
RUN mvn package
CMD ["java", "-jar", "my.jar"]
