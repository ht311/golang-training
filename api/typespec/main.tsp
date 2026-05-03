import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

using TypeSpec.Http;
using TypeSpec.Rest;

@service(#{
  title: "Blog API",
})
@server("http://localhost:8080", "Local development server")
namespace BlogApi;

model Post {
  id: string;
  title: string;
  body: string;
  createdAt: utcDateTime;
  updatedAt: utcDateTime;
}

model CreatePostRequest {
  title: string;
  body: string;
}

model UpdatePostRequest {
  title?: string;
  body?: string;
}

model PostList {
  posts: Post[];
  total: int32;
}

model ErrorResponse {
  message: string;
}

@route("/posts")
interface Posts {
  @get
  list(): PostList | ErrorResponse;

  @post
  create(@body body: CreatePostRequest): {
    @statusCode statusCode: 201;
    @body post: Post;
  } | ErrorResponse;

  @get
  @route("/{id}")
  get(@path id: string): Post | {
    @statusCode statusCode: 404;
    @body error: ErrorResponse;
  };

  @put
  @route("/{id}")
  update(@path id: string, @body body: UpdatePostRequest): Post | {
    @statusCode statusCode: 404;
    @body error: ErrorResponse;
  };

  @delete
  @route("/{id}")
  delete(@path id: string): {
    @statusCode statusCode: 204;
  } | {
    @statusCode statusCode: 404;
    @body error: ErrorResponse;
  };
}

@route("/health")
interface Health {
  @get
  check(): { status: string };
}
