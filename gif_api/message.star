load("render.star", "render")


def main(config):
    return render.Root(
        child=render.Box(
            render.Marquee(
                width=64,
                child=render.Text("local"),
            )
        )
    )
