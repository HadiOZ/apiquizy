FROM ubuntu:latest

WORKDIR /usr/local/api-quizzipy

RUN mkdir /usr/local/data/
RUN mkdir /usr/local/data/profile
RUN mkdir /usr/local/data/quiz
RUN mkdir /usr/local/data/question

ENV PROFILE /usr/local/data/profile
ENV QUIZ /usr/local/data/quiz
ENV QUESTION /usr/local/data/question

COPY apiquizyfull /usr/local/bin

ENV PORT=8000

CMD ["apiquizyfull"]