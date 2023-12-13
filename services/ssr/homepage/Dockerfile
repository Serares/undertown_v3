FROM public.ecr.aws/docker/library/node:18.19.0-slim
COPY --from=public.ecr.aws/awsguru/aws-lambda-adapter:0.7.1 /lambda-adapter /opt/extensions/lambda-adapter
EXPOSE 8080
RUN mkdir /homepage
WORKDIR /homepage
# TODO this looks like a bad idea
COPY ./ ./
RUN yarn
# RUN npx esbuild app.ts --outdir=dist --platform=node --format=cjs --target=node18 --minify
RUN yarn build
RUN cp -rp views ./dist/views
CMD ["node" "dist/app.js"]