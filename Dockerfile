FROM klakegg/hugo:ext-ubuntu

WORKDIR /src

EXPOSE 1313

CMD ["hugo", "server", "--bind", "0.0.0.0"]