background { color rgb <0.1, 0.1, 0.1> }  // Dark gray
#declare FishColor = rgb<0.3, 0.8, 0.8>;
#declare EyeColor = rgb<1, 1, 1>;

#declare Fish = union {
    // Main body - elongated ellipsoid
    sphere {
        <0, 0, 0>, 1
        scale <1.8, 0.8, 0.9>  // Elongated, slightly flattened
        pigment { FishColor }
    }

    // Eyes
    sphere {
        <1.4, 0.2, 0.4>, 0.3
        pigment { EyeColor }
    }
    sphere {
        <1.4, 0.2, -0.4>, 0.3
        pigment { EyeColor }
    }


    // Tail fin
    triangle {
        <-3, 0.9, 0>,
        <-3, -0.9, 0>,
        <-1.5, 0, 0>
        pigment { FishColor }
    }
}

camera {
    location <4, 1.5, -10>
    look_at <0, 0, 0>
    angle 45
}

light_source {
    <-10, 10, -10>
    color rgb 0.7
}
light_source {
    <10, 5, -5>
    color rgb 0.3
}
light_source {
    <0, -5, -10>
    color rgb 0.2
}

object {
    Fish
    rotate <0, 30, 0>
}
