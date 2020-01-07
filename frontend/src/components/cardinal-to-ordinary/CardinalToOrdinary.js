function CardinalToOrdinary(grade) {
    const Ordinary = [
        "Kindergarten",
        "First",
        "Second",
        "Third",
        "Fourth",
        "Fifth",
        "Sixth",
        "Seventh",
        "Eighth",
        "Nineth",
        "Tenth",
        "Eleventh",
        "Twelfth"
    ];

    return Ordinary[grade] + " grade";
}

export default CardinalToOrdinary;
