use actix_web::{get, post, web, App, HttpResponse, HttpServer, Responder};

#[get("/")]
async fn hello() -> impl Responder {
    HttpResponse::Ok().body("Hello world!")
}

#[post("/echo")]
async fn echo(req_body: String) -> impl Responder {
    HttpResponse::Ok().body(req_body)
}

async fn manual_hello() -> impl Responder {
    HttpResponse::Ok().body("Hey there!")
}


#[get("/api/v1/get_google_news")]
async fn get_google_news() -> impl Responder {
    let result = go_browser::get_google_news();
    let json = serde_json::to_string(&result).unwrap();
    HttpResponse::Ok().body(json)
}


pub async fn run_server() -> Result<(),std::io::Error> {
    HttpServer::new(|| {
        App::new()
            .service(hello)
            .service(echo)
            .service(get_google_news)
            .route("/hey", web::get().to(manual_hello))
    })
    .bind(("127.0.0.1", 8080))?
    .run()
    .await
}